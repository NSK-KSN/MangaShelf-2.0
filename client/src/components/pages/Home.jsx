import { loremIpsum } from 'lorem-ipsum';
import { Link } from 'react-router-dom';

const Home = () => (
    <>
      <h3>Home Page</h3>
      <div>
        Home page content: { loremIpsum({ count: 5 })}
      </div>
    </>
  );

  export default Home