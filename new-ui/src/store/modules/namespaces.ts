import { Module } from "vuex";
import { State } from "./../index";
import * as apiNamespace from "../api/namespaces";

export interface NamespacesState {
  namespace: any;
  billInfoData: any;
  billing: any;
  namespaces: Array<any>;
  invoices: Array<any>;
  defaultPerPage: number;
  invoicesLength: number;
  numberNamespaces: number;
  owner: boolean;
}

export const namespaces: Module<NamespacesState, State> = {
  namespaced: true,
  state: {
    namespace: {},
    billInfoData: {},
    billing: {},
    namespaces: [],
    invoices: [],
    defaultPerPage: 3,
    invoicesLength: 0,
    numberNamespaces: 0,
    owner: false,
  },

  getters: {
    list: (state) => state.namespaces,
    get: (state) => state.namespace,
    getNumberNamespaces: (state) => state.numberNamespaces,
    owner: (state) => state.owner,
    billing: (state) => state.billing,
  },

  mutations: {
    setNamespaces: (state, res) => {
      state.namespaces = res.data;
      state.numberNamespaces = parseInt(res.headers["x-total-count"], 10);
    },

    setNamespace: (state, res) => {
      state.namespace = res.data;
    },

    setBilling: (state, data) => {
      state.billing = data;
    },

    removeNamespace: (state, id) => {
      state.namespaces.splice(
        state.namespaces.findIndex((d) => d.tenant_id === id),
        1
      );
    },

    removeMember: (state, usr) => {
      state.namespace.members.splice(
        state.namespace.members.findIndex((m: any) => m.name === usr),
        1
      );
    },

    clearNamespaceList: (state) => {
      state.namespaces = [];
      state.numberNamespaces = 0;
    },

    clearObjectNamespace: (state) => {
      state.namespace = {};
    },

    setOwnerStatus: (state, status) => {
      state.owner = status;
    },
  },

  actions: {
    post: async (context, data) => {
      try {
        const res = await apiNamespace.postNamespace(data);
        return res;
      } catch (error) {
        throw error;
      }
    },

    fetch: async (context, data) => {
      try {
        const res = await apiNamespace.fetchNamespaces(data.page, data.perPage, data.filter);
        context.commit("setNamespaces", res);
      } catch (error) {
        throw error;
      }
    },

    get: async (context, id) => {
      try {
        const res = await apiNamespace.getNamespace(id);
        context.commit("setNamespace", res);
  
        const { billing } = res.data;
        if (billing !== null) {
          context.commit("setBilling", billing);
        }
      } catch (error) {
        throw error;
      }
    },

    put: async (context, data) => {
      try {
        await apiNamespace.putNamespace(data);
      } catch (error) {
        throw error;
      }
    },

    remove: async (context, id) => {
      try {
        await apiNamespace.removeNamespace(id);
        context.commit("removeNamespace", id);
        context.commit("clearObjectNamespace");
        context.commit("clearNamespaceList");
      } catch (error) {
        throw error;
      }
    },

    addUser: async (context, data) => {
      try {
        await apiNamespace.addUserToNamespace(data);
      } catch (error) {
        throw error;
      }
    },

    editUser: async (context, data) => {
      try {
        await apiNamespace.editUserToNamespace(data);
      } catch (error) {
        throw error;
      }
    },

    removeUser: async (context, data) => {
      try {
        await apiNamespace.removeUserFromNamespace(data);
      } catch (error) {
        throw error;
      }
    },

    clearNamespaceList: (context) => {
      context.commit("clearNamespaceList");
    },

    switchNamespace: async (context, data) => {
      try {
        localStorage.removeItem("role");
  
        const res = await apiNamespace.tenantSwitch(data);
        if (res.status === 200) {
          localStorage.setItem("token", res.data.token || "");
          localStorage.setItem("tenant", data.tenant_id);
          localStorage.setItem("role", res.data.role || "");
        }
      } catch (error) {
        throw error;
      }
    },

    setOwnerStatus: async (context, status) => {
      context.commit("setOwnerStatus", status);
    },
  },
};
