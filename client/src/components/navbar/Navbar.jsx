import React from 'react'
import './Navbar.css'
import logo from '../../assets/logo-no-background.png'
import toggle_light_icon from '../../assets/night.png'
import toggle_dark_icon from '../../assets/day.png'
import { Link } from 'react-router-dom';
import { IoMdSearch } from "react-icons/io";

const Navbar = ({theme, setTheme}) => {

    const toggle_mode = () => {
        theme === 'light' ? setTheme('dark') : setTheme('light');
    }

    return (
        <div className='navbar'>
            <Link to="/"><img src={logo} alt="" className='logo' /></Link>
            <ul>
                <li><Link to="/releases">График релизов</Link></li>
                <li><Link to="/browse">Каталог</Link></li>
                <li>Профиль</li>
                <li><Link to="/admin">Администрация</Link></li>
            </ul>
            <div className='search-box'>
                <input type='text' placeholder='Поиск...'/>
                <IoMdSearch className='search-icon'/>
            </div>
            <img onClick={()=>{toggle_mode()}} src={theme === 'light' ? toggle_light_icon : toggle_dark_icon} alt="" className='toggle-icon' />
        </div>
    )
}

export default Navbar