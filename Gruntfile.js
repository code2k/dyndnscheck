'use strict';

module.exports = function (grunt) {
  // load all grunt tasks
  require('load-grunt-tasks')(grunt);

  grunt.initConfig({
    watch: {
      run: {
        files: ['*.go'],
        tasks: ['exec:run'],
        options: {
          spawn: false,
        },
      },
      gruntfile: {
        files: ['Gruntfile.js']
      }
    },
    exec: {
      run: {
        command: 'go run dyndnscheck.go',
      }
    }
  });

  grunt.registerTask('default', [
    'watch'
  ]);
};
