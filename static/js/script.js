window.onload = function() {
	var curveDiv = document.getElementById('Curve');
	var decBtn = document.getElementById('DecDepth');
	var incBtn = document.getElementById('IncDepth');
	var extraForm = document.getElementById('ExtraForm');
	var thicknessRange = document.getElementById('Thickness');
	var colorField = document.getElementById('Color');
	var submitBtn = document.getElementById('Submit');
	var navLinks = document.getElementsByClassName('nav-link');

	var getURL = function(name, depth, onlySVG) {
		var url = '/' + name;
		if (depth !== null) {
			url += '/' + depth;
		}

		params = [];
		if (_thickness !== '3') {
			params.push('t=' + encodeURIComponent(_thickness));
		}
		if (_color !== 'black') {
			params.push('c=' + encodeURIComponent(_color));
		}
		if (onlySVG) {
			params.push('svg=1');
		}
		if (params.length > 0) {
			url += '?' + params.join('&');
		}

		return url;
	}

	var updateURL = function() {
		var url = getURL(_NAME, _depth, false);
		window.history.replaceState(null, document.title, url);
	}

	var updateDecInc = function() {
		if (_depth > 0) {
			decBtn.className = 'depth';
			decBtn.href = getURL(_NAME, _depth - 1, false);
		} else {
			decBtn.className = 'depth disabled';
			decBtn.href = '#';
		}

		if (_depth < _MAX_DEPTH) {
			incBtn.className = 'depth';
			incBtn.href = getURL(_NAME, _depth + 1, false);
		} else {
			incBtn.className = 'depth disabled';
			incBtn.href = '#';
		}
	}

	var updateNavLinks = function() {
		for (var i = 0; i < navLinks.length; i++) {
			link = navLinks[i];
			link.href = getURL(link.innerHTML, null, false);
		}
	}

	var reloadSVG = function() {
		var xhr = new XMLHttpRequest();
		xhr.open('GET', getURL(_NAME, _depth, true), true);
		xhr.onreadystatechange = function() {
			if (xhr.readyState === 4 && xhr.status === 200) {
				curveDiv.innerHTML = xhr.responseText;
			}
		}
		xhr.send();
	}

	var setStyle = function(name, value) {
		var svg = curveDiv.firstElementChild;
		var defs = svg.firstElementChild;
		var style = defs.firstElementChild;
		var rule = style.firstChild;
		var regex = new RegExp(name + ':[^;]+');
		rule.nodeValue = rule.nodeValue.replace(regex, name + ': ' + value);
	}

	submitBtn.style.display = 'none';

	decBtn.addEventListener('click', function(e) {
		e.preventDefault();
		_depth -= 1;
		reloadSVG();
		updateURL();
		updateDecInc();
	});

	incBtn.addEventListener('click', function(e) {
		e.preventDefault();
		_depth += 1;
		reloadSVG();
		updateURL();
		updateDecInc();
	});

	thicknessRange.addEventListener('change', function() {
		_thickness = thicknessRange.value;
		var rect = curveDiv.firstElementChild.viewBox.baseVal;
		var adjusted = _thickness * Math.max(rect.width, rect.height) / 600.0;
		updateURL();
		updateNavLinks();
		setStyle('stroke-width', adjusted);
	});

	extraForm.addEventListener('submit', function(e) {
		e.preventDefault();
		_color = colorField.value;
		updateURL();
		updateNavLinks();
		setStyle('stroke', _color);
	})
}
