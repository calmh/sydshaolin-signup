var editorApp = angular.module('Signup', []);

editorApp.controller('signupCtrl', function ($scope, $location, $http) {
    $scope.form = {dobMonth: "January"};

    $scope.submit = function () {
        $scope.saving = true;
        $http.post("/signup/post", $scope.form).success(function () {
            $scope.saved = true;
        });
    };
});

