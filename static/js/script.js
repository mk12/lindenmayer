function getQuery() {
	var query = {};
	var tokens = window.location.search.substr(1).split('&');
	for (var i = 0; i < tokens.length; i++) {
		if (tokens[i] !== undefined && tokens[i] !== '') {
			var assignment = tokens[i].split('=');
			var name = decodeURIComponent(assignment[0]);
			var value = decodeURIComponent(assignment[1]);
			query[name] = value;
		}
	}
	return query;
}

function setQuery(query) {
	var href = window.location.pathname;
	var i = 0;
	for (var name in query) {
		if (query.hasOwnProperty(name) && query[name]) {
			href += (i === 0) ? '?' : '&';
			var encName = encodeURIComponent(name);
			var encValue = encodeURIComponent(query[name]);
			href += encName + '=' + encValue;
			i++;
		}
	}

	window.location.href = href;
}

window.onload = function() {
	var decBtn = document.getElementById('DecDepth');
	var incBtn = document.getElementById('IncDepth');
	var thicknessRange = document.getElementById('Thickness');
	var colorField = document.getElementById('Color');
	var query = getQuery();

	var load = function() {
		thicknessRange.value = query.s ? query.s : '2';
		colorField.value = query.c ? query.c : '';
	}

	var update = function() {
		query.s = thicknessRange.value;
		query.c = colorField.value;
		setQuery(query);
	}

	load();

	thicknessRange.addEventListener('change', update);
	colorField.addEventListener('keyup', function(e) {
		if (e.keyCode == 13) {
			update();
		}
	});
}
