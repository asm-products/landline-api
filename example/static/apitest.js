var boundLog = console.log.bind(console); // Just a helper to make it easier to pass console.log as a callback.
var Session = function(token){
    this.token = token;
    this.authorizeXHR = function(xhr){
        xhr.setRequestHeader("Authorization", "Bearer "+this.token)
    }

    this.makeCall = function(method, path, data, host){
        host = host || "http://localhost:3000";
        var xhr = new XMLHttpRequest();
        xhr.open(method, host + path);
        this.authorizeXHR(xhr);
         if(data){
            data = JSON.stringify(data);
            xhr.setRequestHeader('Content-Type', 'application/json');
        }
        xhr.send(data);
        return new Promise(function(resolve, reject){
            xhr.addEventListener("load", function(){
                resolve(JSON.parse(xhr.responseText));
            });
            xhr.addEventListener("error", reject);
        });
    }
}

Session.create = function(){
    return new Promise(function(resolve, reject) {
        var xhr = new XMLHttpRequest();
        xhr.open("GET", "http://localhost:3000/sessions/new?team=test-dev");
        xhr.addEventListener("load", function(){
            var resp = JSON.parse(xhr.responseText);
            resolve(new Session(resp.token));
        });
        xhr.addEventListener("error", reject);
        xhr.send();
    });
}

Session.immediate = function(){
    var sess = new Session(null);
    Session.create().then(function(newsess){
        sess.token = newsess.token;
    });
    return sess
}