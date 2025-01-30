import React, { useState, useEffect } from "react";
import axios from "axios";
import "../css/Home.css";
import { FaHeart,FaRegHeart } from "react-icons/fa";
import API_BASE_URL from "../config";
import { useAuth } from "../context/AuthContext";
import { useNavigate } from "react-router-dom";
import { fetchLoggedinUser } from "../utils/api";

function Home() {
    const navigate = useNavigate();
    const { loggedInUser, setLoggedInUser } = useAuth();
    const [isSessionChecked, setIsSessionChecked] = useState(false);
    const [users, setUsers] = useState({});
    const [posts, setPosts] = useState([]);
    const [newPost, setNewPost] = useState({ title: "", description: "" });

    useEffect(() => {
        const verifySession = async () => {
            const user = await fetchLoggedinUser();
            setLoggedInUser(user);
        };
        verifySession();
        setIsSessionChecked(true);
    }, []);

    useEffect(() => {
        if (loggedInUser == null){
            navigate("/login");
        }else{
            fetchData();
            const intervalId = setInterval(fetchData, 10000);

            return () => clearInterval(intervalId);
        }
    }, [isSessionChecked, navigate]);


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
            console.error("Error fetching users: ", error);
        }
    }

    const sortPostsByUpdatedAt = (posts) => {
        posts.sort((a,b) => {
            if(a.updated_at > b.updated_at) return -1;
            else if(a.updated_at < b.updated_at) return 1;

            return 0;
        })
    }

    const fetchPosts = async (updatedUsers) => {
        try {
            const response = await axios.get(`${API_BASE_URL}/posts`);
            const postsData = response.data;

            let postDataWithUserID = postsData.map(postData => ({
                id: postData.id,
                title: postData.title,
                description: postData.description,
                favorites: postData.favorites,
                submitted_by: updatedUsers[postData.submitted_by],
                created_at: postData.created_at,
                updated_at: postData.updated_at,
            }));

            sortPostsByUpdatedAt(postDataWithUserID);
            
            setPosts(postDataWithUserID);
        } catch (error) {
            console.error("Error fetching posts: ", error);
        }
    }

    const fetchData = async () => {
        try{
            const updatedUsers = await fetchUsers();
            await fetchPosts(updatedUsers);
        }catch(error){
            console.error("Error fetching data: ", error);
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
                        <span>
                            <a href="#"><FaRegHeart /></a>
                            <a href="#">{post.favorites.length}</a>
                        </span>
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
