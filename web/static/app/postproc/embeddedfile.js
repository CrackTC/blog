function ReplaceEmbeddedFile() {
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
