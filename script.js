// Copyright 2023 Mitchell Kember. Subject to the MIT License.

// Default state to use when none is given in the URL fragment.
const defaultState = {
    name: "koch",
    depth: 2,
    thickness: 3,
    color: "black",
};

function decodeState(string) {
    if (string === "") return { ...defaultState };
    const parts = string.split("-");
    return {
        name: parts[0] ?? defaultState.name,
        depth: parseInt(parts[1] ?? defaultState.depth),
        thickness: parseFloat(parts[2] ?? defaultState.thickness),
        color: parts[3] ?? defaultState.color,
    };
}

function encodeState(state) {
    return `${state.name}-${state.depth}-${state.thickness}-${state.color}`;
}

// Scaling factor to control the approximate size of the SVG viewBox.
const stepFactor = 600;

// Controls the amount of padding in the viewBox around the curve.
const padFactor = 0.8;

// Lindenmayer systems. Each one has:
// - axiom: initial symbols
// - rules: rules for expanding symbols
// - angle: the angle that "+" and "-" turn, in radians
// - turn: if true, turn the initial direction based on depth
// - base: base of the system size's exponential growth
// - min: minimum depth (corresponds to user depth of 0)
// - max: maximum depth (corresponds to user depth of max-min)
const systems = {
    "koch": {
        axiom: "F++F++F",
        rules: {
            "F": "F-F++F-F",
        },
        angle: Math.PI / 3,
        start: 0,
        turn: false,
        base: 3,
        min: 0,
        max: 7,
    },
    "hilbert": {
        axiom: "a",
        rules: {
            "a": "+bF-aFa-Fb+",
            "b": "-aF+bFb+Fa-",
        },
        angle: Math.PI / 2,
        start: 0,
        turn: false,
        base: 2,
        min: 1,
        max: 8,
    },
    "peano": {
        axiom: "a",
        rules: {
            "a": "aFbFa-F-bFaFb+F+aFbFa",
            "b": "bFaFb+F+aFbFa-F-bFaFb",
        },
        angle: Math.PI / 2,
        start: Math.PI / 2,
        turn: false,
        base: 3,
        min: 1,
        max: 5,
    },
    "gosper": {
        axiom: "A",
        rules: {
            "A": "A-B--B+A++AA+B-",
            "B": "+A-BB--B-A++A+B",
        },
        angle: Math.PI / 3,
        start: Math.PI / 9,
        turn: true,
        base: 2.6,
        min: 0,
        max: 5,
    },
    "sierpinski": {
        axiom: "A",
        rules: {
            "A": "+B-A-B+",
            "B": "-A+B+A-",
        },
        angle: Math.PI / 3,
        start: 0,
        turn: false,
        base: 2,
        min: 1,
        max: 9,
    },
    "rings": {
        axiom: "F+F+F+F",
        rules: {
            "F": "FF+F+F+F+F+F-F",
        },
        angle: Math.PI / 2,
        start: -37 * Math.PI / 360,
        turn: true,
        base: 3,
        min: 0,
        max: 5,
    },
    "tree": {
        axiom: "A",
        rules: {
            "A": "B[+A]-A",
            "B": "BB",
        },
        angle: Math.PI / 4,
        start: Math.PI / 2,
        turn: false,
        base: 1.9,
        min: 0,
        max: 9,
    },
    "plant": {
        axiom: "a",
        rules: {
            "a": "F+[[a]-a]-F[-Fa]+a",
            "F": "FF",
        },
        angle: 25.0 / 180.0 * Math.PI,
        start: Math.PI / 4,
        turn: false,
        base: 2,
        min: 1,
        max: 7,
    },
    "willow": {
        axiom: "a",
        rules: {
            "a": "bFF[+a]c",
            "b": "bF",
            "c": "bFF[-a]a",
        },
        angle: Math.PI / 6,
        start: 80.0 / 180 * Math.PI,
        turn: false,
        base: 1.3,
        min: 1,
        max: 12,
    },
    "dragon": {
        axiom: "Fa",
        rules: {
            "a": "a-bF-",
            "b": "+Fa+b",
        },
        angle: Math.PI / 2,
        start: Math.PI / 4,
        turn: true,
        base: 1.4,
        min: 0,
        max: 15,
    },
    "island": {
        axiom: "F+F+F+F",
        rules: {
            "F": "F+F-F-FF+F+F-F",
        },
        angle: Math.PI / 2,
        start: Math.PI / 4,
        turn: false,
        base: 4,
        min: 0,
        max: 4,
    },
};

// Draws the curve for the given state, returning a string containing SVG.
function draw(state) {
    const system = systems[state.name];
    const depth = system.min + state.depth;
    const step = stepFactor * Math.pow(system.base, -depth);
    const stack = [];
    let t = { x: 0, y: 0, vx: 0, vy: 0, dir: 0 };
    function fmt(f) { return +f.toFixed(3); }
    function coord(k) { return k + fmt(t.x) + " " + fmt(t.y); }
    let d = coord("M");
    let minX = 0, minY = 0, maxX = 0, maxY = 0;
    function advance() {
        t.x += t.vx;
        // Subtract so that origin is in the bottom-left, not top-right.
        t.y -= t.vy;
        if (t.x < minX) minX = t.x;
        if (t.y < minY) minY = t.y;
        if (t.x > maxX) maxX = t.x;
        if (t.y > maxY) maxY = t.y;
        d += coord("L");
    }
    function rotate(r) {
        t.dir += r;
        t.vx = step * Math.cos(t.dir);
        t.vy = step * Math.sin(t.dir);
    }
    function execute(symbol) {
        switch (symbol) {
            case 'F':
            case 'A':
            case 'B':
                advance();
                break;
            case '+':
                rotate(system.angle);
                break;
            case '-':
                rotate(-system.angle);
                break;
            case '[':
                stack.push({ ...t });
                break;
            case ']':
                t = stack.pop();
                d += coord("M");
                break;
        }
    }
    function recur(symbols, depth) {
        if (depth === 0) {
            for (const s of symbols) execute(s);
            return;
        }
        for (const s of symbols) {
            const replacement = system.rules[s];
            if (replacement !== undefined) {
                recur(replacement, depth - 1);
            } else {
                execute(s);
            }
        }
    }
    rotate(system.start * (system.turn ? depth : 1));
    recur(system.axiom, depth);
    const edge = padFactor * state.thickness;
    minX -= edge;
    minY -= edge;
    maxX += edge;
    maxY += edge;
    const width = maxX - minX;
    const height = maxY - minY;
    return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="${minX} ${minY} ${width} ${height}" fill="none" stroke-linecap="square" stroke="${state.color}" stroke-width="${calcStrokeWidth(width, height)}"><path d="${d}"/></svg>`;
}

let state;

const $Curve = document.getElementById("Curve");
const $CurveNameList = document.getElementById("CurveNameList");
const $DecDepth = document.getElementById("DecDepth");
const $IncDepth = document.getElementById("IncDepth");
const $Thickness = document.getElementById("Thickness");
const $Color = document.getElementById("Color");
const $Download = document.getElementById("Download");

function redraw() {
    $Curve.innerHTML = draw(state);
}

function updateUrl() {
    history.replaceState(null, "", "#" + encodeURIComponent(encodeState(state)));
}

function updateDepth() {
    const system = systems[state.name];
    const userMax = system.max - system.min;
    state.depth = Math.max(0, Math.min(userMax, state.depth));
    $DecDepth.classList.toggle("disabled", state.depth === 0);
    $IncDepth.classList.toggle("disabled", state.depth === userMax);
}

let activeCurveLi;
function setActiveCurveLi(a) {
    if (activeCurveLi) {
        activeCurveLi.classList.remove("active");
        activeCurveLi.firstElementChild.classList.remove("disabled");
    }
    activeCurveLi = a;
    activeCurveLi.classList.add("active");
    activeCurveLi.firstElementChild.classList.add("disabled");
}

function viewBoxRect() {
    return $Curve.firstElementChild.viewBox.baseVal;
}

function calcStrokeWidth(width, height) {
    return state.thickness * Math.max(width, height) / stepFactor;
}

function isValidCssColor(string) {
    const style = new Option().style;
    style.color = string;
    return style.color !== "";
}

$DecDepth.addEventListener("click", (e) => {
    e.preventDefault();
    state.depth -= 1;
    updateDepth();
    updateUrl();
    redraw();
});

$IncDepth.addEventListener("click", (e) => {
    e.preventDefault();
    state.depth += 1;
    updateDepth();
    updateUrl();
    redraw();
});

$Thickness.addEventListener("change", () => {
    const delta = $Thickness.value - state.thickness;
    state.thickness = $Thickness.value;
    updateUrl();
    const rect = viewBoxRect();
    const strokeWidth = calcStrokeWidth(rect.width, rect.height);
    $Curve.firstElementChild.setAttribute("stroke-width", strokeWidth);
    const edge = padFactor * delta;
    rect.x -= edge;
    rect.y -= edge;
    rect.width += edge * 2;
    rect.height += edge * 2;
});

$Color.addEventListener("input", () => {
    const color = $Color.value.trim();
    if (!isValidCssColor(color)) return;
    state.color = color;
    updateUrl();
    $Curve.firstElementChild.setAttribute("stroke", state.color);
});

$Download.addEventListener("click", (e) => {
    e.preventDefault();
    const a = document.createElement("a");
    a.setAttribute("href", "data:text/plain;charset=utf-8," + encodeURIComponent(draw(state)));
    a.setAttribute("download", encodeState(state) + ".svg");
    a.style.display = "none";
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
});

for (const li of $CurveNameList.children) {
    const a = li.firstElementChild;
    const name = a.innerHTML;
    a.addEventListener("click", (e) => {
        e.preventDefault();
        state.name = name;
        setActiveCurveLi(li);
        updateDepth();
        updateUrl();
        redraw();
    });
}

function init() {
    state = decodeState(decodeURIComponent(location.hash.slice(1)));
    updateDepth();
    $Thickness.value = state.thickness;
    $Color.value = state.color;
    for (const li of $CurveNameList.children) {
        if (li.firstElementChild.innerHTML === state.name) {
            setActiveCurveLi(li);
        }
    }
    redraw();
}

addEventListener("hashchange", init);
init();
