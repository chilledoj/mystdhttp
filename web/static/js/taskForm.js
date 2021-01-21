const taskId = document.getElementById('taskId').value
const errs = function(el){
  this._el = el
    
  this.render = (err) => {
    const html = err?`<div class="block m-4 p-2"><article class="message is-danger">
    <div class="message-header"><p>Error</p></div>
    <div class="message-body">${err}</div>
    </article></div>`:''
    this._el.innerHTML = html
  } 
  return this
}(document.getElementById('errs'))
const taskForm = document.getElementById('taskForm')

const basePath = '/api/tasks/'
taskForm.addEventListener('submit', function(e){
  e.preventDefault()
 
  let data = Object.fromEntries(new FormData(taskForm))
  errs.render()
  fetch(basePath + (taskId==='new'?'':`${taskId}`), {
    method: taskId==='new'?'POST':'PUT',
    body: JSON.stringify(data),
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
    window.location = meta.headers.get('x-viewlocation')
  })
  .catch(e=>{
    console.log(e)
    errs.render(e.message)
  })
  
});