import React, { useState } from 'react';
import axios from "axios";
import "../css/LoginRegister.css";
import API_BASE_URL from '../config';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

function Login() {
    const [formValue, setFormValue] = useState({
        email: '',
        password: '',
    });

    const {setIsLoggedIn} = useAuth();

    const navigate = useNavigate();

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormValue({
            ...formValue,
            [name]: value,
        });
    };

    const handleSubmit = async(e) => {
        e.preventDefault();
        try {
            let formData = new FormData();
            formData.append("email", formValue.email)
            formData.append("password", formValue.password)

            await axios.post(`${API_BASE_URL}/auth/login`, formData, { withCredentials: true });
            setIsLoggedIn(true);
            navigate("/");
        } catch (error) {
            console.error("Error login:", error);
        }

    };

    return (
        <div>
            <form onSubmit={handleSubmit}>
                <h2>Login</h2>
                <label>
                Email:
                <input
                    type="email"
                    name="email"
                    value={formValue.email}
                    onChange={handleChange}
                    required
                />
                </label>
                <br />
                <label>
                Password:
                <input
                    type="password"
                    name="password"
                    value={formValue.password}
                    onChange={handleChange}
                    required
                />
                </label>
                <br />
                <button type="submit">Login</button>
            </form>
        </div>
    );
}

export default Login;
