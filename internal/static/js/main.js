async function deleteUser(id) {
    if (!confirm('Удалить пользователя?')) return;

    const res = await fetch(`/users/${id}`, { method: 'DELETE' });
    if (res.ok) {
        document.getElementById(`user-${id}`).remove();
    } else {
        alert('Ошибка при удалении пользователя');
    }
}

async function deleteTask(id) {
    if (!confirm('Удалить задачу?')) return;

    const res = await fetch(`/tasks/${id}`, { method: 'DELETE' });
    if (res.ok) {
        document.getElementById(`task-${id}`).remove();
    } else {
        alert('Ошибка при удалении задачи');
    }
}