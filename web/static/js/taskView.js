const form = document.getElementById('commentForm')
form.addEventListener('submit', (e) => {
  e.preventDefault()
  const formData = Object.fromEntries(new FormData(form))
  const taskId = form.dataset.taskId
  console.log(formData, taskId)

  fetch('/api/comments',{
    method: 'POST',
    body: JSON.stringify({ body: formData.cmtBody, taskId }),
    headers: {
      'Content-Type': 'application/json'
    }
  })
  .then(res=>{
    return res.json()
    .then(obj=>({data: obj, meta: res}))
  })
  .then(({data, meta})=>{
    if(!meta.ok){
      throw new Error(`${meta.status} ${meta.statusText}: ${data.error}`)
    }
    form.reset()
    window.location = meta.headers.get('x-viewlocation')
  })
  .catch(e=>{
    console.log(e)
    errs.render(e.message)
  })

})


document.addEventListener('DOMContentLoaded', (event) => {
  const md = new markdownit({typographer: true});
  let tdtls = document.getElementById('taskDetails')
  console.log(tdtls.innerHTML)
  tdtls.innerHTML = `<div class="content">${md.render(tdtls.innerHTML)}</div>`
})