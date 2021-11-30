document.getElementById('data').value  = localStorage.data

document.getElementById('save').onclick = function(){
    localStorage.data = document.getElementById('data').value;
    if (localStorage.data.length == 0) {
        alert('请输入正确的数据')
    } else {
        alert('保存成功.');
    }
}
