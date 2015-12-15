function openD3pieChart(target) {
    'use strict';

    let title, subTitle;

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
        subTitle = "すべて"
    } else {
        subTitle = graphTarget
    }

    $('#myModalLabel').text(title);
    $('#pieChart').text('');

    let pie = new d3pie("pieChart", {
    "header": {
        "title": {
            "text": "",
            "fontSize": 24,
            "font": "open sans"
        },
        "subtitle": {
            "text": "",
            "color": "#999999",
            "fontSize": 12,
            "font": "open sans"
        },
        "titleSubtitlePadding": 9
    },
    "footer": {
        "text": subTitle,
        "color": "#999999",
        "fontSize": 12,
        "font": "open sans",
        "location": "bottom-middle"
    },
    "size": {
        "canvasWidth": 850,
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
    "tooltips": {
        "enabled": true,
            "type": "placeholder",
        "string": "{label}: {value}, {percentage}%"
    },
    "effects": {
        "pullOutSegmentOnClick": {
            "effect": "linear",
            "speed": 400,
            "size": 8
        }
    }
    });


    $('#myModal').modal('show');


}
