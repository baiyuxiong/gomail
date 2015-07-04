var showFlashInfo = function(obj,message){
    var msgDom = $('<div class="alert alert-info" role="alert">'+message+'<div>');
    msgDom.insertBefore(obj)

    window.setTimeout(function(){
        msgDom.fadeOut("slow");
    },2000);
}
var showFlashDanger = function(obj,message){
    var msgDom = $('<div class="alert alert-danger" role="alert">'+message+'<div>');
    msgDom.insertBefore(obj)

    window.setTimeout(function(){
        msgDom.fadeOut("slow");
    },2000);
}