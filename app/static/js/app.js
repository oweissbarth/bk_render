var hub = angular.module("hub",['ngResource', 'ngRoute', 'ngCookies']);

//change template tags to no clash with jinja tags
hub.config(["$interpolateProvider",function($interpolateProvider) {
		$interpolateProvider.startSymbol('[[');
		$interpolateProvider.endSymbol(']]');
		}]);


hub.config(["$locationProvider", "$routeProvider", "$httpProvider", "$cookiesProvider", function($locationProvider, $routeProvider, $httpProvider, $cookiesProvider){
	$routeProvider.when("/", {
			templateUrl : "static/partials/hub.login.html",
			controller: "loginCtrl"
		}).when("/register", {
			templateUrl : "static/partials/hub.register.html",
			controller: "registerCtrl"
		}).when("/admin", {
			templateUrl: "static/partials/hub.admin.dashboard.html",
			controller: "adminDashboardCtrl"
		}).when("/my", {
			templateUrl: "static/partials/hub.user.dashboard.html",
			controller: "userDashboardCtrl"

		}).otherwise({
			redirectTo: "/"
		});

}]);

hub.run(function($rootScope, $http, $cookies){
	var token = $cookies.get("authToken");
	if(token){
		$http.defaults.headers.common.Authorization = "Basic "+btoa(token+":");
		$rootScope.loggedin = true;
	}
});

hub.factory("auth", function(){
	return {
		token : "",
		url: ""
	};
});


hub.controller("loginCtrl", ["$scope", "$http", "$location", "auth", "$cookies", function($scope, $http, $location, auth, $cookies){
	$scope.login = function(){
		$http.get("/login", {headers:{"Authorization": "Basic "+btoa($scope.username+":"+$scope.password)}
			}).then(function(response){
				$cookies.put("authToken", response.data.token);
				$http.defaults.headers.common.Authorization = "Basic "+btoa(response.data.token+":");
				$location.path("/my").replace();
			});
	};
}]);

hub.controller("registerCtrl", ["$scope", "$http","$location", function($scope, $http, $location){
	$scope.register = function(){
		console.log("doing post");
		$http.post("/register", {
			username: $scope.username,
			firstname: $scope.firstname,
			lastname: $scope.lastname,
			password: $scope.password,
			email: $scope.email
		}).then(function(){
			$location.path("/login").replace();
		});

	};

}]);

hub.controller("adminDashboardCtrl", ["$scope", "$resource", function($scope, $resource){
	var Course = $resource("/courses/:courseid", {courseid:"@id"});
	var User = $resource("/users/:userid", {userid:"@id"});
	var Group = $resource("/groups/:groupid", {groupid: "@id"});
	//var Permission = $resource("/permission/:permissionid", ·{permissionid: "@id"});

	$scope.courses = Course.query();
	$scope.users = User.query();
	$scope.groups = Group.query();
	//$scope.permission = Permission.query();

	$scope.createCourse = function(){
		var course = new Course();
		course.name = $scope.newCourseName;
		Course.save(course);
		$scope.courses = Course.query();
	};

	$scope.createGroup= function(){
		var group = new Group();
		group.name = $scope.newGroupName;
		Group.save(group);
		$scope.groups = Group.query();
	};

}]);

hub.controller("userDashboardCtrl", function($scope){


});
