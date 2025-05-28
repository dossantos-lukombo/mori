import { mount } from '@vue/test-utils';
import Search from '../components/Search.vue';
import { createStore } from 'vuex';

// Créer un store Vuex mock
const createVuexStore = () => {
  return createStore({
    state: {
      users: {
        allusers: [
          { id: 1, nickname: 'user1', avatar: 'avatar1.jpg' },
          { id: 2, nickname: 'user2', avatar: 'avatar2.jpg' },
          { id: 3, nickname: 'test3', avatar: 'avatar3.jpg' }
        ]
      },
      groups: {
        allGroups: [
          { id: 1, name: 'group1' },
          { id: 2, name: 'group2' },
          { id: 3, name: 'testgroup' }
        ]
      }
    },
    getters: {
      allUsers: state => state.users.allusers,
      allGroups: state => state.groups.allGroups,
      filterUsers: state => searchquery => {
        if (searchquery === '') return [];
        return state.users.allusers.filter(user => 
          user.nickname.toLowerCase().includes(searchquery.toLowerCase())
        );
      },
      filterGroups: state => searchquery => {
        if (searchquery === '') return [];
        return state.groups.allGroups.filter(group => 
          group.name.toLowerCase().includes(searchquery.toLowerCase())
        );
      }
    },
    actions: {
      getAllUsers: jest.fn(),
      getAllGroups: jest.fn()
    }
  });
};

describe('Search', () => {
  let store;
  const mockRouter = {
    push: jest.fn()
  };

  beforeEach(() => {
    store = createVuexStore();
    jest.clearAllMocks();
  });

  const mountComponent = () => {
    return mount(Search, {
      global: {
        plugins: [store],
        mocks: {
          $router: mockRouter
        }
      }
    });
  };

  it('affiche la barre de recherche', () => {
    const wrapper = mountComponent();
    expect(wrapper.find('input[type="text"]').exists()).toBe(true);
  });

  it('filtre les utilisateurs en fonction de la requête de recherche', async () => {
    const wrapper = mountComponent();
    const input = wrapper.find('input[type="text"]');
    
    // Saisir une requête de recherche
    await input.setValue('user');
    
    // Vérifier que les utilisateurs filtrés sont affichés
    await wrapper.vm.$nextTick();
    const userItems = wrapper.findAll('.item-list li');
    expect(userItems.length).toBe(2); // user1 et user2
  });

  it('navigue vers le profil utilisateur lorsqu\'un utilisateur est cliqué', async () => {
    const wrapper = mountComponent();
    
    // Saisir une requête de recherche pour afficher les utilisateurs
    await wrapper.find('input[type="text"]').setValue('user');
    await wrapper.vm.$nextTick();
    
    // Simuler un clic sur le premier utilisateur
    await wrapper.findAll('.item-list li').at(0).trigger('click');
    
    // Vérifier que la navigation a été appelée avec les bons paramètres
    expect(mockRouter.push).toHaveBeenCalledWith({
      name: 'Profile',
      params: { id: 1 }
    });
  });

  it('navigue vers la page du groupe lorsqu\'un groupe est cliqué', async () => {
    const wrapper = mountComponent();
    
    // Saisir une requête de recherche pour afficher les groupes
    await wrapper.find('input[type="text"]').setValue('group');
    await wrapper.vm.$nextTick();
    
    // Simuler un clic sur le premier groupe
    await wrapper.findAll('.item-list li').at(0).trigger('click');
    
    // Vérifier que la navigation a été appelée avec les bons paramètres
    expect(mockRouter.push).toHaveBeenCalledWith({
      name: 'Group',
      params: { id: 1 }
    });
  });
});