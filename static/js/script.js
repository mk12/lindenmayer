window.onload = function() {
	var curveDiv = document.getElementById('Curve');
	var decBtn = document.getElementById('DecDepth');
	var incBtn = document.getElementById('IncDepth');
	var extraForm = document.getElementById('ExtraForm');
	var thicknessRange = document.getElementById('Thickness');
	var colorField = document.getElementById('Color');
	var submitBtn = document.getElementById('Submit');

	var setStyle = function(name, value) {
		var svg = curveDiv.firstElementChild;
		var defs = svg.firstElementChild;
		var style = defs.firstElementChild;
		var rule = style.firstChild;
		var regex = new RegExp(name + ':[^;]+');
		rule.nodeValue = rule.nodeValue.replace(regex, name + ': ' + value);
	}

	var updateURL = function() {
		var url = '/' + _NAME + '/' + _depth;
		url += '?' + 's=' + _thickness + '&c=' + _color;
		window.history.replaceState(null, document.title, url);
	}

	submitBtn.style.display = 'none';

	thicknessRange.addEventListener('change', function() {
		_thickness = thicknessRange.value;
		setStyle('stroke-width', _thickness);
		updateURL();
	});

	extraForm.addEventListener('submit', function(e) {
		e.preventDefault();
		_color = colorField.value;
		setStyle('stroke', _color);
		updateURL();
	})
}
