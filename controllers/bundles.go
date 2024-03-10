package controllers

import "Mw7/internal/middlewares"

var IndexHandlerGetBundle = middlewares.Join(indexHandlerGet, middlewares.Log, middlewares.UserCheck)
var IndexHandlerPutBundle = middlewares.Join(indexHandlerPut, middlewares.Log, middlewares.Guard)
var IndexHandlerDeleteBundle = middlewares.Join(indexHandlerDelete, middlewares.Log, middlewares.Guard)
var IndexHandlerNoMethBundle = middlewares.Join(indexHandlerNoMeth, middlewares.Log, middlewares.UserCheck)
var IndexHandlerOtherBundle = middlewares.Join(indexHandlerOther, middlewares.Log, middlewares.UserCheck)
var LoginHandlerGetBundle = middlewares.Join(loginHandlerGet, middlewares.Log, middlewares.OnlyVisitors)
var LoginHandlerPostBundle = middlewares.Join(loginHandlerPost, middlewares.Log, middlewares.OnlyVisitors)
var RegisterHandlerGetBundle = middlewares.Join(registerHandlerGet, middlewares.Log, middlewares.OnlyVisitors)
var RegisterHandlerPostBundle = middlewares.Join(registerHandlerPost, middlewares.Log, middlewares.OnlyVisitors)
var HomeHandlerGetBundle = middlewares.Join(homeHandlerGet, middlewares.Log, middlewares.Guard)
var LogHandlerGetBundle = middlewares.Join(logHandlerGet, middlewares.Log, middlewares.Guard) //middlewares.UserCheck)
var ConfirmHandlerGetBundle = middlewares.Join(confirmHandlerGet, middlewares.Log, middlewares.OnlyVisitors)
var ConfirmupdateHandlerGetBundle = middlewares.Join(confirmupdateHandlerGet, middlewares.Log, middlewares.OnlyVisitors)
var LogoutHandlerGetBundle = middlewares.Join(logoutHandlerGet, middlewares.Log, middlewares.Guard)
var ModifUserHandlerGetBundle = middlewares.Join(ModifUserHandlerGet, middlewares.Log, middlewares.Guard)
var ModifUserHandlerPostBundle = middlewares.Join(ModifUserHandlerPost, middlewares.Log, middlewares.Guard)
