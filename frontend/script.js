document.getElementById('addBtn').addEventListener('click', () => {
    const taskInput = document.getElementById('task');
    const task = taskInput.value;

    if (!task) return; // Prevent adding empty tasks

    fetch('/todos', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ task })
    }).then(response => {
        if (response.ok) {
            taskInput.value = ''; // Clear input
            loadTodos(); // Reload the todo list
        } else {
            console.error("Failed to add todo");
        }
    }).catch(err => console.error("Error:", err));
});

function loadTodos() {
    fetch('/todos')
        .then(response => {
            if (!response.ok) throw new Error("Failed to load todos");
            return response.json();
        })
        .then(todos => {
            const todoList = document.getElementById('todoList');
            todoList.innerHTML = '';
            todos.forEach(todo => {
                const li = document.createElement('li');
                li.textContent = todo.task;
                todoList.appendChild(li);
            });
        })
        .catch(err => console.error("Error:", err));
}

// Load todos initially only after login
function onLoginSuccess() {
    document.getElementById('todoContainer').style.display = 'block'; // Show the todo container
    loadTodos(); // Load existing todos
}

// Attach this function to the login event
document.getElementById("loginForm").onsubmit = async function(event) {
    event.preventDefault();
    const username = document.getElementById("loginUsername").value;
    const password = document.getElementById("loginPassword").value;
    const response = await fetch("/login", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ username, password })
    });
    if (response.ok) {
        alert("Logged in successfully!");
        onLoginSuccess(); // Call to load todos on successful login
    } else {
        alert("Invalid credentials. Please try again.");
    }
};
