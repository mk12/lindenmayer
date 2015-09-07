window.onload = function() {
	var curveDiv = document.getElementById('Curve');
	var decBtn = document.getElementById('DecDepth');
	var incBtn = document.getElementById('IncDepth');
	var extraForm = document.getElementById('ExtraForm');
	var thicknessRange = document.getElementById('Thickness');
	var colorField = document.getElementById('Color');
	var submitBtn = document.getElementById('Submit');

	var getURL = function(name, depth, onlySVG) {
		var url = '/' + name + '/' + depth;
		var hasThickness = _thickness !== '3';
		var hasColor = _color !== 'black';

		params = [];
		if (hasThickness) {
			params.push('t=' + encodeURIComponent(_thickness));
		}
		if (hasColor) {
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

	var updateLinks = function() {
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

		// TODO: update navigation links
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
		updateLinks();
	});

	incBtn.addEventListener('click', function(e) {
		e.preventDefault();
		_depth += 1;
		reloadSVG();
		updateURL();
		updateLinks();
	});

	thicknessRange.addEventListener('change', function() {
		_thickness = thicknessRange.value;
		var rect = curveDiv.firstElementChild.viewBox.baseVal;
		var adjusted = _thickness * Math.max(rect.width, rect.height) / 600.0;
		updateURL();
		setStyle('stroke-width', adjusted);
	});

	extraForm.addEventListener('submit', function(e) {
		e.preventDefault();
		_color = colorField.value;
		updateURL();
		setStyle('stroke', _color);
	})
}
