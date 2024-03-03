import React, { useState, useEffect } from 'react';
import './MangaList.css';
import { useNavigate } from 'react-router-dom';

function MangaList() {
    const [data, setData] = useState([]);
    const [search, setSearch] = useState('');
    const navigate = useNavigate();

    useEffect(() => {
        fetch('http://localhost:8080/data') // backend URL
            .then(response => response.json())
            .then(async data => {
                // Sorting the data by score before setting it in the state
                const sortedData = data.sort((a, b) => b.score - a.score);
                setData(sortedData);

                const newData = await Promise.all(data.map(async item => {

                    const typeResponse = await fetch(`http://localhost:8080/dbquery/types/${item.type_id}`);
                    const typeData = await typeResponse.json();

                    return { ...item, type_name: typeData.name };
                }));
                setData(newData);
            })
            .catch(error => console.error('Error fetching data:', error));
    }, []);

    return (
        <div className='container'>
        <div className='filter-bar'>
            <input onChange={(e) => setSearch(e.target.value.toLowerCase())} />
            <select>
                <option value="someOption">Some option</option>
                <option value="otherOption">Other option</option>
            </select>
            <select>
                <option value="someOption">Some option</option>
                <option value="otherOption">Other option</option>
            </select>
            <select>
                <option value="someOption">Some option</option>
                <option value="otherOption">Other option</option>
            </select>
        </div>
        <div className='manga-list'>
            <ul>
                {data.filter((item) => {
                    return search.toLowerCase() === '' ? item : item.title.toLowerCase().includes(search);
                }).map((item) => (
                    <li key={item.id} onClick={() => navigate('/manga/' + item.id)}>
                        <div className='img-container'>
                            <img src={item.cover_image} alt={item.title} />
                        </div>
                        <p>{item.title}</p>
                        <p>Тип: {item.type_name}</p>
                        <p>Рейтинг: {item.score}</p>
                    </li>
                ))}
            </ul>
        </div>
        <h1>Something</h1>
        </div>
    );
}

export default MangaList;
