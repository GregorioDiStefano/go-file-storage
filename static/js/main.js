var app = angular.module('myUpload', ['smoothScroll', 'ngFileUpload', 'ui.bootstrap']);
app.config(function($interpolateProvider) {
    $interpolateProvider.startSymbol('[[');
    $interpolateProvider.endSymbol(']]');
});


app.directive('customOnChange', function() {
    return {
        restrict: 'A',
        link: function(scope, element, attrs) {
            var onChangeHandler = scope.$eval(attrs.customOnChange);
            element.bind('change', onChangeHandler);
        }
    };
});

var newPercent = 0;
app.controller('uploadCtrl', ['$scope', 'Upload', "$http", "$timeout", function($scope, Upload, $http, $timeout) {

    removeProgress = function() {
        setTimeout(function() {
            $scope.$apply(function() {
                newPercent = 0
                $scope.data = {
                    progress: 0,
                    show: false
                };
            });
        }, 200);
    };

    $scope.data = {
        progress: 0,
        show: false
    };

    (function progress() {
            $timeout(function() {
                if (newPercent > 0) {
                    $scope.data.show = true
                    $scope.data.progress = newPercent
                if ($scope.data.progress == 100) {
                    console.log($scope.data.progress)
                    $scope.data.progress = "Syncing..."
                }
              }
                progress();
            }, 50);

    })();

    $scope.uploadFile = function(file) {
        if (file) {
          console.log("Uploading: ", file)
        file.upload = Upload.http({
            method: 'PUT',
            url: location.protocol + '//' + location.hostname + "/" + file.name,
            data: file,
        }).then(function(resp) {
            swal({
                title: "Upload complete!",
                text: "<p style='text-align:left'>Download URL: " + resp.data.downloadURL + "<br> Delete URL: " + resp.data.deleteURL + "</p>",
                customClass: 'swal-wide',
                html: true
            });
            removeProgress()
        }, function(resp) {
            console.log('Error status: ' + resp.status);
            if (resp.status == "429") {
              swal({
                title: "Oops",
                text: "Too many requests, please try again in an hour."
              });
              removeProgress()
            }
        }, function(evt) {
            var progressPercentage = parseInt(100.0 * evt.loaded / evt.total);
            newPercent = progressPercentage
        });
    }
  }
}]);
