var gotrace = angular.module('gotraceback', ['ngCookies']);

gotrace.controller('TracebackCtrl', ['$scope', '$http', '$cookies', function ($scope, $http, $cookies) {
	$scope.analyze = function(traceback) {
		$http({
			method : 'POST',
			url : '/process',
			data : traceback
		}).success(function(data, status, headers, config) {
			// $scope.result = angular.toJson(data);
			$scope.result = data;
		}).error(function(data, status, headers, config) {
		});
	};
	$scope.godocURL = $cookies.godocURL;
	$scope.$watch('godocURL', function(val) {
		$cookies.godocURL = val
	});
}]);
