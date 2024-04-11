interface Project {
  title: string,
  description: string,
  href?: string,
  imgSrc?: string,
}

const projectsData: Project[] = [
  {
    title: 'A Golang API from scratch',
    description: `Learn how to build a RESTful API with Go. This tutorial will walk you through
    the steps of creating a simple API that can be used to create, read, update, and
    delete data from a database.`,
    imgSrc: '/static/images/google.png',
    href: 'https://www.google.com',
  },
  {
    title: 'The Time Machine',
    description: `Imagine being able to travel back in time or to the future. Simple turn the knob
    to the desired date and press "Go". No more worrying about lost keys or
    forgotten headphones with this simple yet affordable solution.`,
    imgSrc: '/static/images/time-machine.jpg',
    href: '/blog/the-time-machine',
  },
]

export default projectsData
