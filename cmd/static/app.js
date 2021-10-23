console.log("JS Loaded")

const url = "127.0.0.1:8080"

function parseOutputOrderJson(outputOrder) {
    col = []
    for(let key in outputOrder) {
        col.push(outputOrder[key]);
    }
    return col
}

var inputForm = document.getElementById("inputForm")
inputForm.addEventListener("submit", (e)=>{
    e.preventDefault()
    const formdata = new FormData(inputForm)
    fetch(document.documentURI,{
        method:"POST",
        body:formdata,
    }).then(
        response => response.json()
    ).then(
        (data) => {
            col = parseOutputOrderJson(data)
            let table = document.getElementById("example")
            let tr = table.insertRow(-1)
            for (let j = 0; j < col.length; j++) {
                let tabCell = tr.insertCell(-1)
                tabCell.innerHTML = col[j]
            }
        }
    ).catch(
        error => console.error(error)
    )
})
