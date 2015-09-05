window.onload = function() {
	var decBtn = document.getElementById('DecDepth');
	var incBtn = document.getElementById('IncDepth');
	var thicknessRange = document.getElementById('Thickness');
	var colorField = document.getElementById('Color');

	var update = function() {
		var path = window.location.pathname;
		path += '?s=' + thicknessRange.value;
		var color = colorField.value.trim();
		if (color !== '') {
			path += '&c=' + color;
		}
		window.location.href = path;
	}

	thicknessRange.addEventListener('change', update);
}
