import { loremIpsum } from 'lorem-ipsum';

const Home = () => (
    <>
      <h3>Home Page</h3>
      <div>
        Home page content: { loremIpsum({ count: 5 })}
      </div>
    </>
  );

  export default Home