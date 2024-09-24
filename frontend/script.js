document.getElementById('addBtn').addEventListener('click', () => {
    const taskInput = document.getElementById('task');
    const task = taskInput.value;

    fetch('http://localhost:3000/todos', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ task })
    }).then(response => response.json())
      .then(data => {
          console.log(data);
          taskInput.value = '';
          loadTodos();
      });
});

function loadTodos() {
    fetch('http://localhost:3000/todos')
        .then(response => response.json())
        .then(todos => {
            const todoList = document.getElementById('todoList');
            todoList.innerHTML = '';
            todos.forEach(todo => {
                const li = document.createElement('li');
                li.textContent = todo.task;
                todoList.appendChild(li);
            });
        });
}

loadTodos();

