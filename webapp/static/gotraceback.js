var gotrace = angular.module('gotraceback', ['ngCookies']);

gotrace.controller('TracebackCtrl', ['$scope', '$http', '$cookies', function ($scope, $http, $cookies) {
	$scope.analyze = function(traceback, format) {
		$http({
			method : 'POST',
			url : '/process?format=' + format,
			data : traceback
		}).success(function(data, status, headers, config) {
			if(format == "json") {
				$scope.result = angular.toJson(data);
			}else{
				$scope.result = data;
			}
		}).error(function(data, status, headers, config) {
		});
	};
	$scope.godocURL = $cookies.godocURL;
	$scope.$watch('godocURL', function(val) {
		$cookies.godocURL = val
	});
}]);
