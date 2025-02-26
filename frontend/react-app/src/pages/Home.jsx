import React, { useState, useEffect, useRef } from "react";
import axios from "axios";
import "../css/Home.css";
import { FaHeart,FaRegHeart } from "react-icons/fa";
import API_BASE_URL from "../config";
import { useAuth } from "../context/AuthContext";
import { useCsrf } from "../context/CsrfContext"
import { useNavigate } from "react-router-dom";
import { fetchLoggedinUser } from "../utils/api";

function Home() {
    const navigate = useNavigate();
    const { loggedInUser, setLoggedInUser } = useAuth();
    const [isSessionChecked, setIsSessionChecked] = useState(false);
    const [users, setUsers] = useState({});
    const [posts, setPosts] = useState([]); // dbから取得してきた投稿
    const [displayedPosts, setDisplayedPosts] = useState([]); // 表示する投稿
    const [page, setPage] = useState(1); // 現在のページ数
    const observer = useRef(null); // intersection Observer の参照
    const POSTS_PER_PAGE = 3; //1ページに表示する投稿数
    const POSTS_PER_FETCH = 300;
    const [newPost, setNewPost] = useState({ title: "", description: "" }); 
    const {csrfToken} = useCsrf();

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

    useEffect(() => {
        setDisplayedPosts(posts.slice(0, POSTS_PER_PAGE));
        setPage(1);
    }, [posts]);
    
    // 下までスクロールしたら新たに表示されるデータが追加される
    const lastPostRef = useRef(null);
    useEffect(() => {
        if (!lastPostRef.current) return;

        observer.current = new IntersectionObserver((entries) => {
            if(entries[0].isIntersecting){
                loadMorePosts();
            }
        },
    { threshold: 1.0 }
    );

        observer.current.observe(lastPostRef.current);
        return () => observer.current.disconnect();
    }, [displayedPosts]);

    // 新しい投稿を追加する関数
    const loadMorePosts = () => {
        const nextPage = page + 1;
        const newPosts = posts.slice(0, nextPage * POSTS_PER_PAGE);

        if(newPosts.length !== displayedPosts.length){
            setDisplayedPosts(newPosts);
            setPage(nextPage);
        }
    };


    const favoriteToggle = async(postId, favoriteAdd) => {
        if(favoriteAdd){
            await axios.post(`${API_BASE_URL}/posts/favorite/${postId}`);
        }else{
            await axios.delete(`${API_BASE_URL}/posts/favorite/${postId}`);
        }

        await fetchData();
    }


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

    const fetchPosts = async (updatedUsers) => {
        try {
            const response = await axios.get(`${API_BASE_URL}/posts`, {params:{limit: POSTS_PER_FETCH}});
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

            await axios.post(`${API_BASE_URL}/posts`, formData, 
                {
                    withCredentials: true,
                    headers: {
                        'X-CSRF-Token': csrfToken
                    }
                }
            );
            setNewPost({ title: "", description: "" });
            fetchData();
        } catch (error) {
            console.error("Error creating post:", error);
        }
    };

    const handleDelete = async (id) => {
        try {
            await axios.delete(`${API_BASE_URL}/posts/${id}`, 
                {
                    withCredentials: true,
                    headers: {
                        'X-CSRF-Token': csrfToken
                    }
                }
            );
            fetchData();
        } catch (error) {
            console.error("Error deleting post:", error);
        }
    };

    const sortPostsByUpdatedAt = (posts) => {
        posts.sort((a,b) => {
            if(a.updated_at > b.updated_at) return -1;
            else if(a.updated_at < b.updated_at) return 1;

            return 0;
        })
    }


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
                    {displayedPosts.map((post, index) => (
                        <li
                            key={post.id}
                            className="post"
                            ref={index === displayedPosts.length - 1 ? lastPostRef : null}
                        >
                            <h3>{post.title}</h3>
                            <span>Submitted by: </span>
                            <a href="#" className="submitted-by">
                                {post.submitted_by?.name || "Anonymous"}
                            </a>
                            <p>{post.description}</p>
                            { post.favorites.find((favoriteUser) => favoriteUser.id == loggedInUser.id) === undefined? (
                                <span>
                                    <button onClick={() => favoriteToggle(post.id, true)}><FaRegHeart /></button>
                                    <a href="#">{post.favorites.length}</a>
                                </span>
                            ):(
                                <span>
                                    <button onClick={() => favoriteToggle(post.id, false)}><FaHeart /></button>
                                    <a href="#">{post.favorites.length}</a>
                                </span>
                            )}
                            <button onClick={() => handleDelete(post.id)} className="deleteBtn">
                                Delete
                            </button>
                        </li>
                    ))}
                </ul>
            </section>
            </main>
        </div>
    );
}

export default Home;
