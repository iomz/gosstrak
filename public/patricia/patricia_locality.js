function loadJSON(callback) {
    var xobj = new XMLHttpRequest();
    xobj.overrideMimeType("application/json");
    xobj.open('GET', 'locality.json', true);
    xobj.onreadystatechange = function () {
        if (xobj.readyState == 4 && xobj.status == "200") {
            callback(xobj.responseText);
        }
    };
    xobj.send(null);
}

function drawTree(treeData) {
    var margin = {
        top: 200,
        right: 120,
        bottom: 20,
        left: 120
    }, width = 960*10000 - margin.right - margin.left, height = 500*10 - margin.top - margin.bottom;

    var i = 0;

    var tree = d3.layout.tree().size([ height, width ]);

    var diagonal = d3.svg.diagonal().projection(function(d) {
        return [ d.x, d.y ];
    });

    var svg = d3.select("body").append("svg").attr("width", width + margin.right + margin.left).attr("height", height + margin.top + margin.bottom).append("g").attr("transform", "translate(" + margin.left + "," + margin.top + ")");

    root = treeData[0];

    update(root, i, tree, diagonal, svg);
}

function getColor(d) {
    var v = d.value;
    if (v > 80) {
        return "#E74C3C";
    } else if (v > 60) {
        return "#EC7063";
    } else if (v > 40) {
        return "#F1948A";
    } else if (v > 20) {
        return "#F5B7B1";
    }
    return "#FADBD8";
}

function update(source, i, tree, diagonal, svg) {
    var nodes = tree.nodes(root).reverse(), links = tree.links(nodes);
    nodes.forEach(function(d) {
        d.y = d.depth * 100;
    });
    var node = svg.selectAll("g.node").data(nodes, function(d) {
        return d.id || (d.id = ++i);
    });
    var nodeEnter = node.enter().append("g").attr("class", "node").attr("transform", function(d) {
        return "translate(" + d.x + "," + d.y + ")";
    });
    nodeEnter.append("circle").attr("r", function(d) { return d.value; }).style("fill", "#FFF");
    nodeEnter.append("text").attr("y", function(d) {
        return d.children || d._children ? -18 : 18;
    }).attr("dy", ".35em").attr("text-anchor", "middle").text(function(d) {
        if (d.value > 20) {
            return d.name;
        }
        return "";
    }).style("fill-opacity", 1);
    var link = svg.selectAll("path.link").data(links, function(d) {
        return d.target.id;
    });
    link.enter().insert("path", "g").attr("class", "link").attr("d", diagonal);
}

loadJSON(function(response) {
    // Parse JSON string into object
    var treeData = JSON.parse(response);
    drawTree(treeData);
});
