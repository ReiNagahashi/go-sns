import React, { useEffect } from "react";
import {Link, useNavigate } from "react-router-dom";
import { checkSession, logoutUser } from "../utils/api";
import { useAuth } from "../context/AuthContext";

const Navbar = () => {
    const { isLoggedIn, setIsLoggedIn } = useAuth();
    const navigate = useNavigate();


    useEffect(() => {
        const verifySession = async () => {
            const loggedIn = await checkSession();
            setIsLoggedIn(loggedIn);
        };
        verifySession();
    }, []);

    const handleLogout = async() => {
        try{
            await logoutUser();
            setIsLoggedIn(false);
            navigate("/login");
        }catch(error){
            console.error("Logout failed: ", error);
        }
    };

    return(
        <nav className="navbar">
            <ul className="navbar-list">
                {isLoggedIn ? (
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