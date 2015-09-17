var gulp = require('gulp'),
    browserify = require('browserify'),
    reactify  = require('reactify'),
    source = require('vinyl-source-stream'),
    inject = require('gulp-inject'),
    path = require('path');
    package  = require('./package.json');

gulp.task('js', function(){
    var extensions = ['.js','.json','.es6'];
    return browserify(package.paths.app)
        .transform(reactify)
        .bundle()
        .pipe(source(package.dest.app))
        .pipe(gulp.dest(package.dest.dist));
}); 

gulp.task('inject', function(){
    return gulp.src("index.html")
            .pipe(
                    inject(
                        gulp.src(package.paths.flexboxcss,{
                            read: false
                        }), {
                            name: "flexboxcss",
                            relative: true
                        })
                 )
           // .pipe(
           //         inject(
           //             gulp.src(package.paths.materialjs,{
           //                 read: false
           //             }), {
           //                 name: "materialjs",
           //                 relative: true
           //             })
           //     )
            .pipe(gulp.dest(package.dest.path))
});
