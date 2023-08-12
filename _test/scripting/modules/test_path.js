var path = require("path");

module.exports = {

    run: function () {
        try {
            var response = {
                delimiter: path.delimiter,
                sep: path.sep,
                dirname: path.dirname("/dir/dir2/dir3/file"),
                basename: path.basename("/dir/dir2/dir3/file.html", ".html"),
                extname: path.extname("/dir/dir2/dir3/file.html"),
                format: path.format({
                    "dir":"/home/dir1",
                    "base":"file.html"
                }),
                isAbsolute: path.isAbsolute("/dir/dir2/dir3/file.html"),
                isAbsolute2: path.isAbsolute("file.html"),
                join: path.join("/dir/dir2", "dir3", "file.html", ".."),
                normalize: path.normalize('/foo/bar//baz/asdf/quux/..'),
                parse: path.parse('/foo/bar/baz/asdf/quux/file.txt'),
                relative: path.relative('/data/orandea/test/aaa', '/data/orandea/impl/bbb'),
                resolve: path.resolve('/foo/bar', './baz'),
                resolve2: path.resolve('/foo/bar', '/tmp/file/'),
                resolve3: path.resolve('wwwroot', 'static_files/png/', '../gif/image.gif'),

            };

            var s = JSON.stringify(response);
            console.debug(s);
            return s;
        } catch (err) {
            console.error("test_path.js", err);
            return err;
        }
    }
}