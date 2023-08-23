(function () {
    function replaceWikiLink() {
        var elements = document.getElementsByClassName("wiki-link")
        for (var i = 0; i < elements.length; i++) {
            let element = elements[i]
            var href = element.getAttribute("href") + ".md"
            fetch("/api/wiki?name=" + encodeURIComponent(href))
                .then(r => r.json())
                .then(data => element.setAttribute("href", data.dest))
        }
    }

    function replaceEmbeddedFile() {
        var elements = document.getElementsByClassName("embedded-file")
        function ReplaceImage(element, name, attr) {
            var path = "/static/blog/" + window.location.pathname.split('/').slice(2, -1).join('/')
            var src = path + '/img/' + encodeURIComponent(name)
            var img = document.createElement("img")
            img.src = src
            img.alt = name
            if (attr) {
                img.width = attr
            }
            element.replaceWith(img)
        }
        for (var i = 0; i < elements.length; i++) {
            var element = elements[i]
            var src = element.getAttribute("data-src")
            var attr = element.getAttribute("data-attr")
            var ext = src.split(".").pop()
            switch (ext) {
                case "md":
                    // markdown file
                    break;
                case "png": case "jpg": case "jpeg": case "gif": case "bmp": case "svg":
                    // image file
                    ReplaceImage(element, src, attr)
                    break;
                case "mp3": case "webm": case "wav": case "m4a": case "ogg": case "3gp": case "flac":
                    // audio file
                    break;
                case "mp4": case "webm": case "ogv":
                    // video file
                    break;
                case "pdf":
                    // pdf file
                    break;
                default:
                    break;
            }
        }
    }

    // css: .table-wrapper { overflow-x: auto; }
    function wrapTableOverflow() {
        var elements = document.getElementsByTagName("table")
        for (var i = 0; i < elements.length; i++) {
            var element = elements[i]
            var wrapper = document.createElement("div")
            wrapper.className = "table-wrapper"
            element.parentElement.insertBefore(wrapper, element)
            wrapper.appendChild(element)
        }
    }

    // css: .code-wrapper { overflow-x: auto; }
    function wrapCodeBlockOverflow() {
        var elements = document.getElementsByTagName("pre")
        for (var i = 0; i < elements.length; i++) {
            var element = elements[i]
            var wrapper = document.createElement("div")
            wrapper.className = "code-wrapper"
            element.parentElement.insertBefore(wrapper, element)
            wrapper.appendChild(element)
        }
    }

    wrapTableOverflow();
    wrapCodeBlockOverflow();
    replaceWikiLink();
    replaceEmbeddedFile();
    hljs.highlightAll();
})();
