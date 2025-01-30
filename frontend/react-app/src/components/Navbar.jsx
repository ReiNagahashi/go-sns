import React, { useEffect } from "react";
import {Link, useNavigate } from "react-router-dom";
import { fetchLoggedinUser, logoutUser } from "../utils/api";
import { useAuth } from "../context/AuthContext";

const Navbar = () => {
    const { loggedInUser, setLoggedInUser } = useAuth();
    const navigate = useNavigate();

    useEffect(() => {
        const verifySession = async () => {
            const user = await fetchLoggedinUser();
            setLoggedInUser(user);
        };
        verifySession();
    }, []);

    const handleLogout = async() => {
        try{
            await logoutUser();
            setLoggedInUser(null);
            navigate("/login");
        }catch(error){
            console.error("Logout failed: ", error);
        }
    };

    return(
        <nav className="navbar">
            <ul className="navbar-list">
                {loggedInUser != null ? (
                    <>
                        <li>
                            <Link to="/">Home</Link>
                        </li>
                        <li>
                            <button onClick={handleLogout} className="logout-button">
                                Logout
                            </button>
                        </li>
                    </>
                ) : (
                    <>
                        <li>
                            <Link to="/login">Login</Link>
                        </li>
                        <li>
                            <Link to="/register">Register</Link>
                        </li>
                    </>
                )}
            </ul>
        </nav>
    );
};

export default Navbar;