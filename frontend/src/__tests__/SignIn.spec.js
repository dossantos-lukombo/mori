import { mount } from '@vue/test-utils';
import SignIn from '../components/SignIn.vue';
import { createApp } from 'vue';

// Configuration minimale de Vue Router pour les tests
const mockRouter = {
  push: jest.fn(),
  currentRoute: { value: { path: '/' } }
};

// Mock the store dispatch method
const mockStore = {
  dispatch: jest.fn().mockResolvedValue()
};

// Create a stub for router-link component
const RouterLinkStub = {
  name: 'RouterLink',
  props: ['to'],
  template: '<a :href="to"><slot /></a>'
};

describe('SignIn', () => {
  const globalConfig = {
    global: {
      plugins: [],
      mocks: {
        $router: mockRouter,
        $toast: {
          open: jest.fn()
        },
        $store: mockStore
      },
      stubs: {
        'router-link': RouterLinkStub
      }
    }
  };

  it('affiche correctement les éléments du formulaire', () => {
    const wrapper = mount(SignIn, globalConfig);
    expect(wrapper.find('form').exists()).toBe(true);
    expect(wrapper.find('#email').exists()).toBe(true);
    expect(wrapper.find('#password').exists()).toBe(true);
  });
});