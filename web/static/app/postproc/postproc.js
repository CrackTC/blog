function PostProcess() {
    ReplaceWikiLink();
    ReplaceEmbeddedFile();
    MathJax.typeset();
}
