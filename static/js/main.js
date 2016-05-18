var app = angular.module('myUpload', ['ngFileUpload']);

app.directive('customOnChange', function() {
  return {
    restrict: 'A',
    link: function (scope, element, attrs) {
      var onChangeHandler = scope.$eval(attrs.customOnChange);
      element.bind('change', onChangeHandler);
    }
  };
});


app.controller('uploadCtrl', ['$scope', 'Upload', function($scope, Upload, $http) {

  $scope.uploadFile = function(file) {
      file.upload = Upload.http({
        method: 'PUT',
        url: 'http://' + location.hostname + "/" + file.name,
        data: file,
      });
    }
}]);
