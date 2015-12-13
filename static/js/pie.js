'use strict';

// Get query parameters
let queries = location.search.replace('?', '').split('&'),
    query = {},
    title = "",
    subTitle = "",
    i = 0;
for (i; i < queries.length; i++) {
    let v = queries[i].split("=");
    query[v[0]] = v[1];
}
let target = query['target']

// Title
let graphType = target.substring(0, 1)
if (graphType === 'b') {
    title = "ブラウザ利用比率"
} else if (graphType == 'o') {
    title = "OS利用比率"
}

// Subtitle
let graphTarget = target.substring(2)
if (graphTarget === 'all') {
    subTitle = "全体"
} else {
    subTitle = graphTarget
}

let pie = new d3pie("pieChart", {
	"header": {
		"title": {
			"text": title,
			"fontSize": 24,
			"font": "open sans"
		},
		"subtitle": {
		    "text": subTitle,
			"color": "#999999",
			"fontSize": 12,
			"font": "open sans"
		},
		"titleSubtitlePadding": 9
	},
	"footer": {
		"color": "#999999",
		"fontSize": 10,
		"font": "open sans",
		"location": "bottom-left"
	},
	"size": {
		"canvasWidth": 650,
		"pieOuterRadius": "90%"
	},
	"data": data_all[target],
	"labels": {
		"outer": {
			"pieDistance": 32
		},
		"inner": {
			"hideWhenLessThanPercentage": 3
		},
		"mainLabel": {
			"fontSize": 11
		},
		"percentage": {
			"color": "#ffffff",
			"decimalPlaces": 0
		},
		"value": {
			"color": "#adadad",
			"fontSize": 11
		},
		"lines": {
			"enabled": true
		},
		"truncation": {
			"enabled": true
		}
	},
	"effects": {
		"pullOutSegmentOnClick": {
			"effect": "linear",
			"speed": 400,
			"size": 8
		}
	}
});
