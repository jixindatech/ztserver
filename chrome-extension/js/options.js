document.getElementById('email').value = localStorage.email
document.getElementById('token').value  = localStorage.token

document.getElementById('save').onclick = function(){
    localStorage.email = document.getElementById('email').value;
    localStorage.token = document.getElementById('token').value;
    if (localStorage.email.length == 0 || localStorage.token == 0) {
        alert('请输入正确的邮箱地址或Token')
    } else {
        alert('保存成功.');
    }
}
