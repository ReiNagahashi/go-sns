import React, { useState, useEffect } from "react";
import axios from "axios";
import "../css/Home.css";
import API_BASE_URL from "../config";

function Home() {
    const [posts, setPosts] = useState([]);
    const [newPost, setNewPost] = useState({ title: "", description: "" });

    // Fetch all posts
    useEffect(() => {
        fetchPosts();
    }, []);

    const fetchPosts = async () => {
        try {
            const response = await axios.get(`${API_BASE_URL}/posts`);
            setPosts(response.data);
        } catch (error) {
            console.error("Error fetching posts:", error);

        }
    };

    const handleInputChange = (e) => {
        const { name, value } = e.target;
        setNewPost({ ...newPost, [name]: value });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            let formData = new FormData();
            formData.append("title", newPost.title)
            formData.append("description", newPost.description)

            await axios.post(API_BASE_URL, formData);
            setNewPost({ title: "", description: "" });
            fetchPosts();
        } catch (error) {
            console.error("Error creating post:", error);
        }
    };

    const handleDelete = async (id) => {
        try {
            await axios.delete(`${API_BASE_URL}/${id}`);
            fetchPosts();
        } catch (error) {
            console.error("Error deleting post:", error);
        }
    };

    return (
        <div className="App">
            <header>
            <h1>Green SNS</h1>
            </header>
            <main>
            <section className="post-form">
                <form onSubmit={handleSubmit}>
                <h2>Create a New Post</h2>
                <input
                    type="text"
                    name="title"
                    placeholder="Title"
                    value={newPost.title}
                    onChange={handleInputChange}
                    required
                />
                <textarea
                    name="description"
                    placeholder="Description"
                    value={newPost.description}
                    onChange={handleInputChange}
                    required
                ></textarea>
                <button type="submit">Post</button>
                </form>
            </section>
            <section className="post-list">
                <h2>Posts</h2>
                <ul>
                {posts.map((post) => (
                    <li key={post.id} className="post">
                    <h3>{post.title}</h3>
                    <p>{post.description}</p>
                    <button onClick={() => handleDelete(post.id)}>Delete</button>
                    </li>
                ))}
                </ul>
            </section>
            </main>
        </div>
    );
}

export default Home;
