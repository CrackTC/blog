function ReplaceWikiLink() {
    var elements = document.getElementsByClassName("wiki-link")
    for (var i = 0; i < elements.length; i++) {
        let element = elements[i]
        var href = element.getAttribute("href") + ".md"
        fetch("/api/wiki?name=" + encodeURIComponent(href))
            .then(r => r.json())
            .then(data => element.setAttribute("href", data.dest))
    }
}
