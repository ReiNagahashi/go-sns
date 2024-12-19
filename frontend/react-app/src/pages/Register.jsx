import React, { useState } from 'react';
import "../css/LoginRegister.css";

function Register() {
    const [formData, setFormData] = useState({
        username: '',
        email: '',
        password: '',
    });

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData({
            ...formData,
            [name]: value,
        });
    };

    const handleSubmit = (e) => {
        e.preventDefault();
        console.log('Registering with data:', formData);
        // API呼び出しを実装
    };

    return (
        <div>
            <form onSubmit={handleSubmit}>
            <h2>Register</h2>
            <label>
                Username:
                <input
                type="text"
                name="username"
                value={formData.username}
                onChange={handleChange}
                required
                />
            </label>
            <br />
            <label>
                Email:
                <input
                type="email"
                name="email"
                value={formData.email}
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
                value={formData.password}
                onChange={handleChange}
                required
                />
            </label>
            <br />
            <button type="submit">Register</button>
            </form>
        </div>
    );
}

export default Register;
