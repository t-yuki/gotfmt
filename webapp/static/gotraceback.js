var gotrace = angular.module('gotraceback', []);

gotrace.controller('TracebackCtrl', ['$scope', '$http', function ($scope, $http) {
	$scope.analyze = function(traceback) {
		$http({
			method : 'POST',
			url : '/process',
			data : traceback
		}).success(function(data, status, headers, config) {
			$scope.result = angular.toJson(data);
		}).error(function(data, status, headers, config) {
		});
	};
}]);
