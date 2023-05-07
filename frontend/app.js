const API_URL = "http://localhost:8080/books";

async function loadBooks() {
    const container = document.getElementById("books");
    container.innerHTML = "";

    const res = await fetch(API_URL);
    const books = await res.json();

    books.forEach(book => {
        const card = document.createElement("div");
        card.className = "book-card";

        card.innerHTML = `
            <p class="book-title">${book.id}. ${book.title}</p>
            <p class="book-author">by ${book.author}</p>

            <button class="btn-small btn-edit" onclick="editBook(${book.id}, '${book.title}', '${book.author}')">
                Edit
            </button>

            <button class="btn-small btn-delete" onclick="deleteBook(${book.id})">
                Delete
            </button>
        `;

        container.appendChild(card);
    });
}

async function addBook() {
    const title = document.getElementById("title").value.trim();
    const author = document.getElementById("author").value.trim();

    if (!title || !author) {
        alert("Both fields are required!");
        return;
    }

    await fetch(API_URL, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ title, author })
    });

    document.getElementById("title").value = "";
    document.getElementById("author").value = "";

    loadBooks();
}

async function deleteBook(id) {
    await fetch(`${API_URL}/${id}`, {
        method: "DELETE"
    });

    loadBooks();
}

async function editBook(id, oldTitle, oldAuthor) {
    const title = prompt("New Title:", oldTitle);
    const author = prompt("New Author:", oldAuthor);

    if (!title || !author) return;

    await fetch(`${API_URL}/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ title, author })
    });

    loadBooks();
}

loadBooks();
