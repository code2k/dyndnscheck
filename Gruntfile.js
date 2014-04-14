'use strict';

module.exports = function (grunt) {
  // load all grunt tasks
  require('load-grunt-tasks')(grunt);

  grunt.initConfig({
    watch: {
      run: {
        files: ['*.go'],
        tasks: ['go:run:osx'],
        options: {
          spawn: false,
        },
      },
      gruntfile: {
        files: ['Gruntfile.js']
      }
    },
    go: {
      osx: {
        output: 'dyndnscheck',
        env: {
          GOARCH: 'amd64',
          GOOS: 'darwin'
        },
        run_files: ['dyndnscheck.go']
      },
      linux: {
        output: 'dyndnscheck',
        env: {
          GOARCH: '386',
          GOOS: 'linux'
        },
        run_files: ['dyndnscheck.go']
      },
      windows: {
        output: 'dyndnscheck.exe',
        env: {
          GOARCH: '386',
          GOOS: 'windows'
        },
        run_files: ['dyndnscheck.go']
      }
    }
  });

  grunt.registerTask('default', [
    'watch'
  ]);

  grunt.registerTask('build', [
    'go:build:osx'
  ]);
};
