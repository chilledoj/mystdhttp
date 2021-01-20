let md

document.addEventListener('DOMContentLoaded', (event) => {
  md = new markdownit();
  reloadData()
})

const taskEl = document.getElementById("tasks");
const delMod = document.getElementById('deleteModal');
const delModBody = document.getElementById('deleteModalBody');

const closeModal = () => {
  delMod.classList.remove('is-active')
}

const showModal = (id, title) => {
  delModBody.innerHTML = `<p class="heading">Are you sure you want to delete "${title}" ?</p>`
  delMod.classList.add("is-active")
  delMod.dataset.deleteId = id
}

const deleteTask = (id) => {
  return fetch(`/api/tasks/${id}`, {
      method: 'DELETE',
    })
    .then(res=>{
      if(!res.ok){
        return res.json()
        .then(obj=>({data: obj, meta: res}))
        .then(({data, meta}) => {
          if(!meta.ok){
            throw new Error(`${meta.status} ${meta.statusText}: ${data.error}`)
          }
          return
        })
      }
      return
    })  
}

document.addEventListener('click', (e) => {
  const deleteLink = e.target.closest('.card-footer-item')
  if(deleteLink && deleteLink.dataset.deleteId){
    e.preventDefault()
    // console.log("delete", deleteLink.dataset.deleteId)
    showModal(deleteLink.dataset.deleteId,deleteLink.dataset.deleteTitle)
    return
  }
  const modal = e.target.closest('.modal')
  if(modal){
    if(e.target.closest('[aria-label="close"]')){
      closeModal()
      return
    }
    if(e.target.dataset.delete){
      deleteTask(modal.dataset.deleteId)
      .then(()=>{
        closeModal()
        reloadData()
      })
      .catch(e=>{
        delModBody.appendChild(document.createTextNode(e.message))
      })
    }
  }
})

const renderStatusTag = (status) => {
  let col = "primary"
  if(status == "inactive"){
    col = "light"
  } else if (status == "started") {
    col = "info"
  } else if (status == "completed") {
    col = "success"
  }
  return `<span class="tag is-${col}">${status}</span>`
}
const twodigit = (n) => ('0' + n.toString()).slice(-2)
const formatDateText = (dateText) => {
  let d
  try{
    d = new Date(dateText)
  }catch(e){
    console.log("Couldn't parse date : %s", dataText)
  }
  return `${d.getFullYear()}-${twodigit(d.getMonth()+1)}-${twodigit(d.getDate())} ${twodigit(d.getHours())}:${twodigit(d.getMinutes())}`
}

const renderTaskCard = (task) => `<div class="column is-one-third">
<div class="card">
    <div class="card-content">
      <div style="float: right;">${renderStatusTag(task.status)}</div>
      <a class="title ${task.status=="inactive"?'has-text-grey-light':'has-text-link'}" href="/view/${task.id}">${task.title}</a>
      <div class="block mt-4 p-4 has-background-white-ter ${task.status=="inactive"?'has-text-grey-light':''}">
        <div class="content">${md.render(task.details.length>100?task.details.substr(0,100):task.details)}</div>
      </div>
      <div class="is-flex is-justify-content-space-between mt-4">
        <span>Assigned to: ${task.assignedTo}</span>
        <span class="subtitle is-size-7 ${task.status=="inactive"?'has-text-grey-light':''}">Last updated: ${formatDateText(task.updAt)}</span>
      </div>
    </div>
    <footer class="card-footer">
      <a href="/view/${task.id}" class="card-footer-item">View</a>
      <a href="/edit/${task.id}" class="card-footer-item">Edit</a>
      <a href="/delete/${task.id}" class="card-footer-item has-text-danger" data-delete-id="${task.id}" data-delete-title="${task.title}" >
        <span class="icon"><i class="gg-trash"></i></span>
        <span class="ml-2"> Delete</span>
      </a>
    </footer>
</div>
</div>`

const reloadData = () => {
  fetch('/api/tasks')
  .then(res=>res.json())
  .then(data=>{
    taskEl.innerHTML = data.sort((a,b)=>{
      return a.title.localeCompare(b.title)
    }).map(renderTaskCard).join('')
  })
  .catch(e=>{
    taskEl.innerHTML = `<div class="column is-full">${e.message}</div>`
  })
}

