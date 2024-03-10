import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import '../../styles/ReleaseSchedule.css';

const ReleaseSchedule = () => {
    const [bookData, setBookData] = useState([]);
    const navigate = useNavigate();
    const [yearFilt, setYearFilt] = useState(2024);

    useEffect(() => {
      fetchData();
    }, []);
  
    const fetchData = async () => {
      try {
        const response = await axios.get('http://localhost:8080/releases-data');
        setBookData(response.data);
      } catch (error) {
        console.error('Error fetching book data:', error);
      }
    };

    const handleYearChange = (year) => {
        setYearFilt(year);
      };
  
      const renderBooksByMonth = () => {
        // Group books by month and year
        const booksByMonth = {};
    
        // Filter books for the selected year
        const booksForYear = bookData.filter(book => new Date(book.release_date).getFullYear() === yearFilt);
    
        booksForYear.forEach((book) => {
          const releaseDate = new Date(book.release_date);
          const monthYearKey = `${releaseDate.getFullYear()}-${(releaseDate.getMonth() + 1).toString().padStart(2, '0')}`;
    
          if (!booksByMonth[monthYearKey]) {
            booksByMonth[monthYearKey] = [];
          }
    
          booksByMonth[monthYearKey].push(book);
        });
    
        // Sort the keys (month-year) in descending order
        const sortedKeys = Object.keys(booksByMonth).sort().reverse();
    
        return sortedKeys.map((key) => {
          const [year, month] = key.split('-');
          const releaseMonth = new Date(`${year}-${month}-01`).toLocaleDateString('en-US', { month: 'long' });
    
          return (
            <div key={key} className='release-page-container'>
              <h2>{releaseMonth} {year}</h2>
              <div className='list-of-manga'>
                {booksByMonth[key].map((book) => (
                  <div className='manga-container' key={book.id} onClick={() => navigate('/manga/' + book.title_id)}>
                    <img src={book.cover_link} alt={book.title} />
                    <p>{book.title}</p>
                    <p>Том {book.volume_number}</p>
                  </div>
                ))}
              </div>
            </div>
          );
        });
    };
    
  
    return (
      <div>
        {Array.from({ length: 9 }, (_, i) => (
          <button key={2016 + i} onClick={() => handleYearChange(2016 + i)}>
            {2016 + i}
          </button>
        ))}
        {renderBooksByMonth()}
      </div>
    );
  }


export default ReleaseSchedule;
