var app = angular.module('myUpload', ['ngFileUpload', 'ui.bootstrap']);
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
    $scope.data = {
        progress: 0,
        show: false
    };

    (function progress() {
        if ($scope.data.progress < 100) {
            $timeout(function() {
                if (newPercent > 0) {
                    $scope.data.progress = newPercent
                    $scope.data.show = true
                }
                progress();
            }, 200);
        }
    })();

    $scope.uploadFile = function(file) {
        file.upload = Upload.http({
            method: 'PUT',
            url: 'http://' + location.hostname + "/" + file.name,
            data: file,
        }).then(function(resp) {
            console.log('Success ' + resp.config.file.name + 'uploaded. Response: ' + resp.data);
        }, function(resp) {
            console.log('Error status: ' + resp.status);
        }, function(evt) {
            var progressPercentage = parseInt(100.0 * evt.loaded / evt.total);
            newPercent = progressPercentage
        });
    }
}]);
