let button = document.querySelector("#submit");
if (button) {
    button.onclick = function (e) {
        let inputs = document.querySelectorAll(".form-control");
        // document.querySelector("#floatingInput")
        
        let data = {};
        for (let i = 0; i < inputs.length; i++) {
            data[inputs[i].name] = inputs[i].value;
        }
        let xhr = new XMLHttpRequest();
        xhr.open("POST", "/user/reg");
        xhr.onload = function (e) {
            let response = JSON.parse(e.currentTarget.response);
            if ("Error" in response) {
                if (response.Error == null) {
                    console.log("Пользователь успешно зарегистрирован");
                    window.location.url = "http://127.0.0.1:8080/";
                } else {
                    console.log(response.Error);
                }
            } else {
                console.log("Некорректные данные");
            }
        };
        xhr.send(JSON.stringify(data));
    }
}
