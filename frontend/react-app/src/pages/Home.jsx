import React, { useState, useEffect } from "react";
import axios from "axios";
import "../css/Home.css";
import API_BASE_URL from "../config";
import { useAuth } from "../context/AuthContext";
import { useNavigate } from "react-router-dom";
import { checkSession } from "../utils/api";

function Home() {
    const navigate = useNavigate();
    const { isLoggedIn, setIsLoggedIn } = useAuth();
    const [isSessionChecked, setIsSessionChecked] = useState(false);
    const [users, setUsers] = useState({});
    const [posts, setPosts] = useState([]);
    const [newPost, setNewPost] = useState({ title: "", description: "" });

    useEffect(() => {
        const verifySession = async () => {
            const loggedIn = await checkSession();

            setIsLoggedIn(loggedIn);
            setIsSessionChecked(true);
        };

        verifySession();
    }, [isLoggedIn]);

    useEffect(() => {
        if (!isSessionChecked) return;

        if (!isLoggedIn){
            navigate("/login");
        }else{
            fetchData();
        }
    }, [isSessionChecked, isLoggedIn, navigate]);


    const fetchUsers = async() => {
        try {
            const response = await axios.get(`${API_BASE_URL}/users`);
            const usersData = response.data;

            // idをキーにユーザーデータを値に入れる
            return usersData.reduce((acc, userData) => {
                acc[userData.id] = userData;
                return acc;
            }, { ...users });

        } catch (error) {
            console.error("Error fetching users:", error);
        }
    }


    const fetchPosts = async (updatedUsers) => {
        try {
            const response = await axios.get(`${API_BASE_URL}/posts`);
            const postsData = response.data;

            const postDataWithUserID = postsData.map(postData => ({
                id: postData.id,
                title: postData.title,
                description: postData.description,
                submitted_by: updatedUsers[postData.submitted_by]
            }));
            
            setPosts(postDataWithUserID);
        } catch (error) {
            console.error("Error fetching posts:", error);
        }
    }
// TODO: mapでposts配列を展開した時に、submitted_by=idとなるuserの情報を表示する
// →O(N)でやるためには、usersデータは配列ではなくディクショナリとして持っておき、users[submitted_by]という感じで展開すればいいと思う。
    const fetchData = async () => {
        try{
            const updatedUsers = await fetchUsers();
            await fetchPosts(updatedUsers);
        }catch(error){
            console.log("Error fetching data: ", error);
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

            await axios.post(`${API_BASE_URL}/posts`, formData);
            setNewPost({ title: "", description: "" });
            fetchData();
        } catch (error) {
            console.error("Error creating post:", error);
        }
    };

    const handleDelete = async (id) => {
        try {
            console.log(`${API_BASE_URL}/posts/${id}`);
            await axios.delete(`${API_BASE_URL}/posts/${id}`);
            fetchData();
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
                        <span>Submitted by: </span>
                        <a href="#" className="submitted-by">
                            {post.submitted_by?.name || "Anonymous"}
                        </a>
                        <p>{post.description}</p>
                        <button onClick={() => handleDelete(post.id)} className="deleteBtn">Delete</button>
                    </li>
                ))}
                </ul>
            </section>
            </main>
        </div>
    );
}

export default Home;
