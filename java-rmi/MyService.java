package com.dkprojects.rmi;
import java.rmi.*;


public interface MyService extends Remote{
    void put(String key, String value) throws RemoteException;
    int get(String key) throws RemoteException;
}