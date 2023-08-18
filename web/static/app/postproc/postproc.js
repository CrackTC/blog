function PostProcess() {
    ReplaceWikiLink();
    ReplaceEmbeddedFile();
    MathJax.typeset();
    hljs.highlightAll();
}
